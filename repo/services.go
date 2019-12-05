package repo

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/makeos/mosdef/types"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing/format/pktline"
)

// service describes a git service and its handler
type service struct {
	method string
	handle func(*serviceParams) error
}

type serviceParams struct {
	w          http.ResponseWriter
	r          *http.Request
	hook       *PushHook
	repo       types.BareRepo
	repoDir    string
	op         string
	srvName    string
	gitBinPath string
}

// getInfoRefs Handle incoming request for a repository's references
func getInfoRefs(s *serviceParams) error {

	var err error
	var refs []byte
	var version string
	var args = []string{s.srvName, "--stateless-rpc", "--advertise-refs", "."}
	var isDumb = s.srvName == ""

	// If this is a request from a dumb client, skip to dumb response section
	if isDumb {
		goto dumbReq
	}

	// Execute git command which will return references
	refs, err = execGitCmd(s.gitBinPath, s.repoDir, args...)
	if err != nil {
		return err
	}

	// Configure response headers. Disable cache and set code to 200
	hdrNoCache(s.w)
	s.w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", s.srvName))
	s.w.WriteHeader(http.StatusOK)

	version = getVersion(s.r)

	// If request is not a protocol v2 request, write the smart parameters
	// describing the service response
	if version != "2" {
		s.w.Write(packetWrite("# service=git-" + s.srvName + "\n"))
		s.w.Write(packetFlush())
	}

	// Write the references received from the git-upload-pack command
	s.w.Write(refs)
	return nil

	// Handle dumb request
dumbReq:

	hdrNoCache(s.w)

	// At this point, the dumb client needs help getting files (since it does
	// not support pack generation on-the-fly). Generate auxiliary files to help
	// the client discover the references and packs the server has.
	updateServerInfo(s.gitBinPath, s.repoDir)

	// Send the info/refs file back to the client
	return sendFile(s.op, "text/plain; charset=utf-8", s)
}

// serveService handles git-upload & fetch-pack requests
func serveService(s *serviceParams) error {
	w, r, op, dir := s.w, s.r, s.op, s.repoDir
	op = strings.ReplaceAll(op, "git-", "")

	// Set response headers
	w.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-result", op))
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	hdrNoCache(w)

	// Construct the git command
	env := os.Environ()
	args := []string{op, "--stateless-rpc", dir}
	cmd := exec.Command(s.gitBinPath, args...)
	version := r.Header.Get("Git-Protocol")
	cmd.Dir = dir
	cmd.Env = env

	// If client requested v2 protocol, set protocol flag in env
	if len(version) != 0 {
		cmd.Env = append(env, fmt.Sprintf("GIT_PROTOCOL=%s", version))
	}

	// Get the command's stdin pipe
	in, err := cmd.StdinPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdin pipe")
	}

	// Get the command's stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}
	defer stdout.Close()

	// Start running the command (does not wait)
	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start command")
	}

	// If the request is compressed, we need to uncompress
	// before we feed it to the git.
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(r.Body)
		defer reader.Close()
	default:
		reader = r.Body
		defer reader.Close()
	}

	// Handle fetch request
	if op == "upload-pack" {
		io.Copy(in, reader)
		in.Close()
		io.Copy(w, stdout)
		return nil
	}

	if err := s.hook.BeforePush(); err != nil {
		return errors.Wrap(err, "BeforePush error")
	}

	// Inspect the pushed data and extract useful information.
	// When we done inspecting, send the pushed data to git.
	var pushReader *PushReader
	pushReader, err = newPushReader(in, s.repo)
	if err != nil {
		return errors.Wrap(err, "unable to create push inspector")
	}
	io.Copy(pushReader, reader)
	if err = pushReader.Read(); err != nil {
		return errors.Wrap(err, "push inspection failed")
	}

	scn := pktline.NewScanner(stdout)
	pktEnc := pktline.NewEncoder(w)
	defer pktEnc.Flush()

	// Do work that needs to be done after git finished processing the pushed data.
	if err := s.hook.AfterPush(pushReader); err != nil {
		pktEnc.Encode(sidebandErr(err.Error()))
		return errors.Wrap(err, "BeforeOutput hook err")
	}

	// Write output from git to the http response
	for scn.Scan() {
		pktEnc.Encode(scn.Bytes())
	}

	return nil
}