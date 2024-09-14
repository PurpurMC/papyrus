package org.purpurmc.papyrus.db.entity;

import jakarta.annotation.Nonnull;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;
import jakarta.persistence.UniqueConstraint;
import java.util.UUID;

@Entity
@Table(uniqueConstraints = @UniqueConstraint(name = "UniqueMetadataNameAndProject", columnNames = {"NAME", "PROJECT_ID"}))
public class ProjectMetadata {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Nonnull
    private String name;

    @Nonnull
    @Column(name = "P_VALUE")
    private String value;

    @Nonnull
    @ManyToOne
    @JoinColumn(name = "PROJECT_ID", referencedColumnName = "ID")
    private Project project;

    public ProjectMetadata(Project project, String name, String value) {
        this.project = project;
        this.name = name;
        this.value = value;
    }

    public ProjectMetadata() {
    }

    public String getName() {
        return name;
    }

    public String getValue() {
        return value;
    }
}
