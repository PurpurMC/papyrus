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
@Table(uniqueConstraints = @UniqueConstraint(name = "UniqueNameAndBuild", columnNames = {"NAME", "BUILD_ID"}))
public class Metadata {
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
    @JoinColumn(name = "BUILD_ID", referencedColumnName = "ID")
    private Build build;

    public Metadata(Build build, String name, String value) {
        this.build = build;
        this.name = name;
        this.value = value;
    }

    public Metadata() {
    }

    public String getName() {
        return name;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }
}
