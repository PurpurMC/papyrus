package org.purpurmc.papyrus;

import org.purpurmc.papyrus.commands.APICommand;
import picocli.CommandLine;

@CommandLine.Command(
        name = "papyrus",
        subcommands = {
                CommandLine.HelpCommand.class,
                APICommand.class,
        }
)
public class Papyrus implements Runnable {

    @Override
    public void run() {
        System.out.println("Papyrus is running!");
    }
}
