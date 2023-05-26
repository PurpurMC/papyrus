package org.purpurmc.papyrus.db.repository;

import org.purpurmc.papyrus.db.entity.CreationState;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;
import java.util.UUID;

public interface CreationStateRepository extends JpaRepository<CreationState, UUID> {
    Optional<CreationState> getStateById(UUID uuid);
}
