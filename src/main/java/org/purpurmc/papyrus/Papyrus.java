package org.purpurmc.papyrus;

import org.purpurmc.papyrus.commands.PapyrusCommand;
import picocli.CommandLine;

public class Papyrus {

    public static void main(String[] args) {
        try {
            if (PapyrusConfig.needsSetup()) {
                PapyrusConfig.setup();
            } else {
                PapyrusConfig.load(false);
            }
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        }

        new CommandLine(new PapyrusCommand()).execute(args);
    }
}
