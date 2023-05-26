package org.purpurmc.papyrus.db.entity;

import jakarta.annotation.Nonnull;
import jakarta.annotation.Nullable;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.OneToOne;
import jakarta.persistence.Table;
import jakarta.persistence.UniqueConstraint;

import java.util.UUID;

@Entity
@Table(uniqueConstraints = @UniqueConstraint(name = "UniqueNameAndVersion", columnNames = {"NAME", "VERSION_ID"}))
public class Build {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Nonnull
    @Column(unique = true)
    private String name;

    @Nonnull
    @Enumerated(EnumType.STRING)
    private BuildResult result;

    @Nonnull
    private Long timestamp;

    @Nonnull
    private Long duration;

    @Nullable
    private String hash;

    @Nonnull
    @ManyToOne
    @JoinColumn(name = "VERSION_ID", referencedColumnName = "ID")
    private Version version;

    @Nullable
    @OneToOne(mappedBy = "build")
    private File file;

    public Build() {
    }

    public Build(Version version, String name, BuildResult result, Long timestamp, Long duration) {
        this(version, name, result, timestamp, duration, null);
    }

    public Build(Version version, String name, BuildResult result, Long timestamp, Long duration, String hash) {
        this.version = version;
        this.name = name;
        this.result = result;
        this.timestamp = timestamp;
        this.duration = duration;
        this.hash = hash;
    }

    public UUID getId() {
        return this.id;
    }

    public String getName() {
        return this.name;
    }

    public BuildResult getResult() {
        return this.result;
    }

    public Long getTimestamp() {
        return this.timestamp;
    }

    public Long getDuration() {
        return this.duration;
    }

    public String getHash() {
        return this.hash;
    }

    public void setHash(String hash) {
        this.hash = hash;
    }

    public enum BuildResult {
        SUCCESS,
        FAILURE
    }
}
