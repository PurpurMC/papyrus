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

    @Query("SELECT b FROM Build b WHERE b.version = :version AND b.name = :name AND b.ready = 1")
    Optional<Build> findByVersionAndNameAndReady(Version version, String name);

    @Query("SELECT b FROM Build b WHERE b.version = :version AND b.ready = 1 ORDER BY timestamp ASC")
    List<Build> findAllByVersionAndReadyOrderByTimestampAsc(Version version);

    @Query("SELECT b FROM Build b WHERE b.version = :version AND b.name = :name AND b.ready = 1 AND b.result = 'SUCCESS'")
    Optional<Build> findByVersionAndNameAndReadyAndResultIsSuccess(Version version, String name);

    @Query("SELECT b FROM Build b WHERE b.version = :version AND b.name = :name AND b.file IS NOT null AND b.result = 'SUCCESS'")
    Optional<Build> findByVersionAndNameAndFileNotNullAndResultIsSuccess(Version version, String name);

    @Query("SELECT b FROM Build b WHERE b.version = :version AND b.file IS NOT null AND b.result = 'SUCCESS' ORDER BY timestamp DESC LIMIT 1")
    Optional<Build> findLatestByVersionAndFileNotNull(Version version);
}
