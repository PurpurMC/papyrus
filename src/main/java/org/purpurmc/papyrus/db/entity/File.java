package org.purpurmc.papyrus.db.entity;

import jakarta.annotation.Nonnull;
import jakarta.annotation.Nullable;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.OneToOne;

import java.util.UUID;

@Entity
public class File {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Nonnull
    private String contentType;

    @Nullable
    private String fileExtension;

    @Nonnull
    @OneToOne
    @JoinColumn(name = "BUILD_ID", referencedColumnName = "ID")
    private Build build;

    public File() {
    }

    public File(Build build, String contentType, String fileExtension) {
        this.build = build;
        this.contentType = contentType;
        this.fileExtension = fileExtension;
    }

    public UUID getId() {
        return this.id;
    }

    public String getContentType() {
        return this.contentType;
    }

    public String getFileExtension() {
        return this.fileExtension;
    }
}
