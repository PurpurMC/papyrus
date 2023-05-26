package org.purpurmc.papyrus.db.repository;

import org.purpurmc.papyrus.db.entity.Build;
import org.purpurmc.papyrus.db.entity.Commit;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;

@Repository
public interface CommitRepository extends JpaRepository<Commit, UUID> {
    List<Commit> findAllByBuild(Build build);
}
