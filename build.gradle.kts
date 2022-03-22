plugins {
    java
    application
}

group = "org.purpurmc"
version = "2.0.0-INDEV"

application {
    mainClass.set("org.purpurmc.papyrus.Main")
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("info.picocli", "picocli", "4.6.3")
    implementation("io.javalin", "javalin", "4.4.0")
    implementation("org.slf4j", "slf4j-simple", "1.7.36")
}
