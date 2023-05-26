package org.purpurmc.papyrus.db.entity;

import jakarta.annotation.Nonnull;
import jakarta.annotation.Nullable;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToOne;

import java.util.Optional;
import java.util.UUID;

@Entity
public class CreationState {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Nullable
    private String fileExtension;

    @Nonnull
    @OneToOne
    @JoinColumn(name = "BUILD_ID", referencedColumnName = "ID")
    private Build build;

    public CreationState() {
    }

    public CreationState(Build build, Optional<String> fileExtension) {
        this.build = build;
        this.fileExtension = fileExtension.orElse(null);
    }

    public UUID getId() {
        return this.id;
    }

    public String getFileExtension() {
        return this.fileExtension;
    }

    public Build getBuild() {
        return this.build;
    }
}
