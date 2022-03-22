package org.purpurmc.papyrus;

import picocli.CommandLine;

public class Main {

    public static void main(String[] args) {
        new CommandLine(new Papyrus()).execute(args);
    }
}
