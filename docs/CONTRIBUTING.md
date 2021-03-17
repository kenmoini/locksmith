# Contributing

> **WARNING:** This guide is a work in progress and will continue to evolve over
> time. If you have content to contribute, please refer to this document
> each time as things may have changed since the last time you contributed.
>
> This warning will be removed once we have settled on a reasonable set of
> guidelines for contributions.

### Before Contributing

Any open-source project is eager for contributions from the community - making the most effective use of everyone's efforts will help making your contributions and having them accepted into the main code base.

Before hacking away at an enhancment of patch, open an Issue and propose your contribution idea.  There are plenty of grand ideas out there, but not all of them are appropriate for the scope of Locksmith and may be better served as another separate project that can use Locksmith as a foundation to quickly innovate upon.

Starting a dialog with the community can provide guidance, help shape ideas into capabilities, and provide visiblity into the efforts of others at large for better collaboration.  All this saves everyone's time and provides you more agility in creating and contributing.

With that being said, this is how you would technically contribute to this project:

### 1. Fork the Locksmith repo

Forking Locksmith is a simple process - you could [click this link to fork the repo](https://github.com/kenmoini/locksmith/fork) or do the following:

1. On GitHub, navigate to the https://github.com/kenmoini/locksmith repo.
2. In the top-right corner of the page, click **Fork**.

That's it! Now, you have a [Fork][git-fork] of the original kenmoini/locksmith repo.

### 2. Create a local clone of your fork

Right now, you have a fork of the Locksmith repo in your GitHub account, but you don't have the files in that repo on your computer. Let's create a [clone][git-clone] of your fork locally on your computer.

```sh
git clone git@github.com:YOUR-USERNAME/locksmith.git
cd locksmith

# Configure git to sync your fork with the original repo
git remote add upstream https://github.com/kenmoini/locksmith

# Never push to upstream repo
git remote set-url --push upstream no_push
```

### 3. Verify your [remotes][git-remotes]

To verify the new upstream repo you've specified for your fork, type `git remote -v`. You should see the URL for your fork as `origin`, and the URL for the original repo as `upstream`.

```sh
origin  git@github.com:YOUR-USERNAME/locksmith.git (fetch)
origin  git@github.com:YOUR-USERNAME/locksmith.git (push)
upstream        https://github.com/kenmoini/locksmith (fetch)
upstream        no_push (push)
```

### 4. Create a new branch

Creating a new branch allows for easier merging, avoiding of conflicts, and better origination analysis across the lifecycle of a code base.  Let's assume you're creating the `awesome-new-feature` branch:

```sh
git branch -M awesome-new-feature
```

You are now using the `awesome-new-feature` branch on your local computer.

### 5. Modify the code base

Now you would implement the Awesome New Feature or a Documentation Fix.  Before submitting your work ensure it passes build requirements with:

```sh
go build
```

Along the way make sure to [Add][git-add] and [Commit][git-commit] your changes - the more atomic and targeted you make your commits the easier it is to be validated and approved.

Say you just modified the README.md file to add some documentation around your Awesome New Feature, as well as an `awesome-new-feature.go` file for the application code base.  This is how you would add and commit it:

```sh
git add README.md
git commit -m "Added Awesome New Feature documentation to readme"

git add awesome-new-feature.go
git commit -m "Added Awesome New Feature go script"
```

At least try to target your commits per file/directory where it makes sense.  Additional detail can be provided in the Pull Request.

### 6. [Push][git-push] your `awesome-new-feature`

When ready to review (or just to establish an offsite backup of your work in your GitHub repo), push your branch to your fork on the GitHub remote server:

```sh
git push -u origin awesome-new-feature
```

A benefit to pushing to your fork on GitHub is that it will automatically run the GitHub Actions to perform a Test of the Build - click the "Actions" tab of your GitHub repo!

### 7. Submit a [pull request][pr]

1. Visit your fork at https://github.com/YOUR-USERNAME/locksmith
2. Use the Branch/Tag drop-down to select your `awesome-new-feature` branch.
3. Click the `Compare` button to the right of your `awesome-new-feature` branch.
4. Click the "Create pull request" button and fill out the PR with as much supporting information as you can and submit.

At this point you're waiting on a Locksmith maintainer.  Some changes or improvements may be suggested or alternatives presented, expect a back and forth of communication before it's integrated into the code base.

[git-fork]: https://help.github.com/articles/fork-a-repo/
[git-clone]: https://git-scm.com/docs/git-clone
[git-remotes]: https://git-scm.com/book/en/v2/Git-Basics-Working-with-Remotes
[git-branch]: https://git-scm.com/docs/git-branch
[git-commit]: https://git-scm.com/docs/git-commit
[git-push]: https://git-scm.com/docs/git-push
[git-add]: https://git-scm.com/docs/git-add
[pr]: https://github.com/kenmoini/locksmith/compare/