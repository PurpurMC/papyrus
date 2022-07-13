package org.purpurmc.papyrus.jenkins.uploader.payload;

import java.util.List;

public class CreateBuildPayload {
    private String project;
    private String version;
    private String build;
    private String result;
    private List<Commit> commits;
    private long duration;
    private long timestamp;

    public void setProject(String project) {
        this.project = project;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public void setBuild(String build) {
        this.build = build;
    }

    public void setResult(String result) {
        this.result = result;
    }

    public void setCommits(List<Commit> commits) {
        this.commits = commits;
    }

    public void setDuration(long duration) {
        this.duration = duration;
    }

    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }

    public static class Commit {
        private String author;
        private String email;
        private String description;
        private String hash;
        private long timestamp;

        public void setAuthor(String author) {
            this.author = author;
        }

        public void setEmail(String email) {
            this.email = email;
        }

        public void setDescription(String description) {
            this.description = description;
        }

        public void setHash(String hash) {
            this.hash = hash;
        }

        public void setTimestamp(long timestamp) {
            this.timestamp = timestamp;
        }
    }
}
