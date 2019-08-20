node("jenkins-ec2") {

    git_commit_short_sha = repoPull(
        sha: GIT_SHA
    )

    app = imageBuild(
        name: 'vault',
        repo_name: 'vault',
        git_short_sha: git_commit_short_sha,
    )

    imagePush(
        image: app
    )

    app = imageBuild(
        name: 'sha_printer',
        repo_name: 'sha_printer',
        git_short_sha: git_commit_short_sha,
        context_dir: 'sha_printer'
    )

    imagePush(
        image: app
    )

}
