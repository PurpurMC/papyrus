package org.purpurmc.papyrus.util;

import org.purpurmc.papyrus.config.AppConfiguration;
import org.purpurmc.papyrus.exception.InvalidAuthToken;

public class AuthUtil {
    public static void requireAuth(AppConfiguration configuration, String authHeader) {
        String[] parts = authHeader.trim().split(" ");
        if (parts.length != 2) {
            throw new InvalidAuthToken();
        }

        if (!parts[0].equals("Basic")) {
            throw new InvalidAuthToken();
        }

        if (!parts[1].equals(configuration.getAuthToken())) {
            throw new InvalidAuthToken();
        }
    }
}
