package org.purpurmc.papyrus.db.entity;

import jakarta.annotation.Nonnull;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;

import java.util.UUID;

@Entity
public class Commit {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Nonnull
    private String author;

    @Nonnull
    private String email;

    @Nonnull
    private String description;

    @Nonnull
    private String hash;

    @Nonnull
    private Long timestamp;

    @Nonnull
    @ManyToOne
    @JoinColumn(name = "BUILD_ID", referencedColumnName = "ID")
    private Build build;

    public Commit() {
    }

    public Commit(Build build, String author, String email, String description, String hash, Long timestamp) {
        this.build = build;
        this.author = author;
        this.email = email;
        this.description = description;
        this.hash = hash;
        this.timestamp = timestamp;
    }

    public String getAuthor() {
        return this.author;
    }

    public String getEmail() {
        return this.email;
    }

    public String getDescription() {
        return this.description;
    }

    public String getHash() {
        return this.hash;
    }

    public Long getTimestamp() {
        return this.timestamp;
    }
}
