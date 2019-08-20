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

}
