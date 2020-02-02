/*
requires libgit2
*/


package main

import (
    "github.com/libgit2/git2go"
    "errors"
    "flag"
)

func main(){

    repoPath := flag.String("repo", "/usr/src/app/libgit2", "path to the git repository")
    flag.Parse()

    repo, err := git.OpenRepository(*repoPath)
    if err != nil {
        panic(err)
    }

    reference, err := repo.References.Lookup("refs/remotes/origin/master")
    if err != nil {
        panic(err)
    }

    branch := reference.Branch()

    // Get the name
    name, err := branch.Name()
    if err != nil {
        panic(err)
    }

    // Locate remote
    remote, err := repo.Remotes.Lookup("origin")
    if err != nil {
        panic(err)
    }

    // Fetch changes from remote
    err = remote.Fetch([]string{}, nil, "")
    if err != nil {
        panic(err)
    }

    // Get remote master
    remoteBranch, err := repo.References.Lookup("refs/remotes/"+name)
    if err != nil {
        panic(err)
    }

    remoteBranchID := remoteBranch.Target()
    // Get annotated commit
    annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
    if err != nil {
        panic(err)
    }

    // Do the merge analysis
    mergeHeads := make([]*git.AnnotatedCommit, 1)
    mergeHeads[0] = annotatedCommit
    analysis, _, err := repo.MergeAnalysis(mergeHeads)
    if err != nil {
        panic(err)
    }

    // Get repo head
    head, err := repo.Head()
    if err != nil {
        panic(err)
    }

    if analysis & git.MergeAnalysisUpToDate != 0 {
        return
    } else if analysis & git.MergeAnalysisNormal != 0 {
        // Just merge changes
        err = repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil)
        if err != nil {
            panic(err)
        }
        // Check for conflicts
        index, err := repo.Index()
        if err != nil {
            panic(err)
        }

        if index.HasConflicts() {
            panic(errors.New("Conflicts encountered. Please resolve them."))
        }

        // Make the merge commit
        sig, err := repo.DefaultSignature()
        if err != nil {
            panic(err)
        }

        // Get Write Tree
        treeId, err := index.WriteTree()
        if err != nil {
            panic(err)
        }

        tree, err := repo.LookupTree(treeId)
        if err != nil {
            panic(err)
        }

        localCommit, err := repo.LookupCommit(head.Target())
        if err != nil {
            panic(err)
        }

        remoteCommit, err := repo.LookupCommit(remoteBranchID)
        if err != nil {
            panic(err)
        }

        repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

        // Clean up
        repo.StateCleanup()
    } else if analysis & git.MergeAnalysisFastForward != 0 {
        // Fast-forward changes
        // Get remote tree
        remoteTree, err := repo.LookupTree(remoteBranchID)
        if err != nil {
            panic(err)
        }

        // Checkout
        err = repo.CheckoutTree(remoteTree, nil)
        if err != nil {
            panic(err)
        }

        branchRef, err := repo.References.Lookup("refs/heads/"+name)
        if err != nil {
            panic(err)
        }

        // Point branch to the object
        branchRef.SetTarget(remoteBranchID, "")
        _, err = head.SetTarget(remoteBranchID, "")
        if err != nil {
            panic(err)
        }

    } else {
        panic(analysis)
    }

}