package gitprocess

import (
	"fmt"

	"com.mooregreatsoftware/go-git-process/vendor/_nuts/github.com/libgit2/git2go"
)

// FetchOptions specifies what to do for the fetch
type FetchOptions struct {
	// RemoteName is the name of the remote URL, such as "orgin"
	RemoteName string

	// Prune will remove branches that no longer exist in the remote from the tracking branches
	Prune bool
}

// Fetch updates the local repository with the contents of a remote repository.
// If there are no remotes then this is a no-op.
func Fetch(repo Repository, options FetchOptions, remoteFactory RemoteFactory, remoteListFactory RemotesFactory) error {
	var remoteName string
	if options.RemoteName == "" {
		remoteNames, err := remoteListFactory()
		if err != nil {
			log.Warning("Could not get remote names")
			return err
		}
		remoteName = normalizeRemoteName(options.RemoteName, remoteNames)
	}

	if remoteName == "" {
		log.Info("This repo has no remotes")
		return nil
	}

	remote, err := remoteFactory(remoteName)
	if err != nil {
		log.Warning("Could not resolve remote for \"%s\"", remoteName)
		return err
	}

	if err := remote.Fetch(); err != nil {
		log.Error("Could not fetch")
		return err
	}
	log.Info("Fetched from \"%s\"", remoteName)

	if options.Prune {
		log.Info("Pruning %s", remote.Name())
		return remote.Prune()
	}

	return nil
}

// ******************************************
//
// Types
//
// ******************************************

// Remotes is a collection of Remote instances
type Remotes []Remote

// RemoteFactory defines how to create an instance of Remote given a name
type RemoteFactory func(remoteName string) (Remote, error)

// RemotesFactory defines how to get to complete set of Remotes available
type RemotesFactory func() (Remotes, error)

// Remote provides ways to connect to a remote git repository
type Remote interface {
	Prune() error // remove obsolete remote branch references
	Fetch() error // retrieve the latest information from the remote repository
	Name() string // the name of the remote repository
}

// ******************************************
//
// gitRemoteStruct
// implements the Remote interface
//
// ******************************************

/*
gitRemoteStruct implements the Remote interface using the git2go library
*/
type gitRemoteStruct struct {
	Repo      Repository
	name      string
	gitRemote *git.Remote
	gitRepo   *git.Repository
}

// Prune removes obsolete remote branch references
func (remote gitRemoteStruct) Prune() error {
	remote.ensureGitRemote()
	return remote.gitRemote.Prune()
}

// Name gives the name of the remote repository
func (remote gitRemoteStruct) Name() string {
	return remote.name
}

// Fetch retrieves the latest information from the remote repository
func (remote gitRemoteStruct) Fetch() error {
	remote.ensureGitRemote()
	log.Info("Fetching %s", remote.Name)
	if err := remote.gitRemote.Fetch(nil, nil, ""); err != nil {
		log.Error(fmt.Sprintf("Could not fetch %s", remote.Name()))
		return err
	}
	return nil
}

func (remote *gitRemoteStruct) ensureGitRepo() {
	if remote.gitRepo == nil {
		gitRepo, err := git.OpenRepository(remote.Repo.Path())
		if err != nil {
			log.Fatalf("Could not open git repo: %s - \"%s\"", remote.Repo.Path(), err)
		}
		remote.gitRepo = gitRepo
	}
}

func (remote *gitRemoteStruct) ensureGitRemote() {
	remote.ensureGitRepo()
	if remote.gitRemote == nil {
		gitRemote, err := gitRemote(remote)
		if err != nil {
			log.Fatalf("Could not get git remote: %s - \"%s\"", remote.Name(), err)
		}
		remote.gitRemote = gitRemote
	}
}

// ******************************************
//
// Remotes
//
// ******************************************

// ContainsName says if one of the Remotes has the given name
func (remotes Remotes) ContainsName(remoteName string) bool {
	return remotes.GetByName(remoteName) != nil
}

// GetByName retrieves the Remote with the given name, or nil if not found
func (remotes Remotes) GetByName(remoteName string) *Remote {
	for _, remote := range remotes {
		if remoteName == remote.Name() {
			return &remote
		}
	}
	return nil
}

// ******************************************
// ******************************************

func listRemotes(gitRepo *git.Repository, repo *Repository, remoteFactory *RemoteFactory) (Remotes, error) {
	remoteNames, err := gitRepo.ListRemotes()
	if err != nil {
		return nil, err
	}

	var remotes Remotes
	for _, remoteName := range remoteNames {
		remote, err := (*remoteFactory)(remoteName)
		if err != nil {
			return nil, err
		}
		remotes = append(remotes, remote)
	}
	return remotes, nil
}

func gitRemote(remote *gitRemoteStruct) (*git.Remote, error) {
	remote.ensureGitRepo()
	remotes, err := remoteNames(remote.gitRepo, remote.Name())
	if err != nil {
		return nil, err
	}
	if len(remotes) == 0 {
		log.Debug("This repo has no remotes")
		return nil, nil
	}
	log.Debug("remotes: %v", remotes)

	log.Info("Looking up %s in %s", remote.Name, remote.gitRepo.Workdir())
	gitRemote, err := remote.gitRepo.LookupRemote(remote.Name())
	if err != nil {
		log.Warning("Could not lookup remote \"%s\"", remote.Name)
		return nil, err
	}

	return gitRemote, nil
}

// normalizeRemoteName tries to make sure that if the remote name is empty
// then a valid one is given back
func normalizeRemoteName(remoteName string, remotes Remotes) string {
	if remoteName == "" {
		if remotes.ContainsName("origin") {
			log.Info("Defaulting remoteName to \"origin\"")
			return "origin"
		}
		// log.Info("Defaulting remoteName to \"\"")
		return "" // TODO: Compute a default value
	}
	return remoteName
}

func remoteNames(gitRepo *git.Repository, remoteName string) ([]string, error) {
	remotes, err := gitRepo.ListRemotes()
	if err != nil {
		return nil, err
	}
	log.Debug("remotes: %v", remotes)

	if len(remotes) == 0 && remoteName != "" {
		return nil, fmt.Errorf("There are no remotes, but \"%s\" was passed in", remoteName)
	}
	return remotes, nil
}
