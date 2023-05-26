package org.purpurmc.papyrus.db.repository;

import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Version;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface BuildRepository extends JpaRepository<Build, UUID> {
    boolean existsByVersionAndName(Version version, String name);

    Optional<Build> findByVersionAndNameAndFileNotNull(Version version, String name);

    List<Build> findAllByVersionAndFileNotNullOrderByTimestampAsc(Version version);

    @Query("SELECT b FROM Build b WHERE b.result = 'SUCCESS' AND b.version = :version AND b.file IS NOT null ORDER BY timestamp DESC LIMIT 1")
    Optional<Build> findLatestByVersionAndFileNotNull(@Param("version") Version version);
}
