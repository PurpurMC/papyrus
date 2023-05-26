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
@Table(uniqueConstraints = @UniqueConstraint(name = "UniqueNameAndProject", columnNames = {"NAME", "PROJECT_ID"}))
public class Version {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Nonnull
    @Column(unique = true)
    private String name;

    @Nonnull
    @ManyToOne
    @JoinColumn(name = "PROJECT_ID", referencedColumnName = "ID")
    private Project project;

    public Version() {
    }

    public Version(Project project, String name) {
        this.project = project;
        this.name = name;
    }

    public String getName() {
        return this.name;
    }
}
