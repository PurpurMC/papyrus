package org.purpurmc.papyrus.db.repository;

import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Metadata;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface MetadataRepository extends JpaRepository<Metadata, UUID> {

    List<Metadata> findByBuild(Build build);
}
