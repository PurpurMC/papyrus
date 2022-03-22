package org.purpurmc.papyrus.commands;

import picocli.CommandLine;

@CommandLine.Command(
        name = "papyrus",
        subcommands = {
                CommandLine.HelpCommand.class,
                APICommand.class,
        }
)
public class PapyrusCommand implements Runnable {

    @Override
    public void run() {
        System.out.println("Papyrus is running!");
    }
}
