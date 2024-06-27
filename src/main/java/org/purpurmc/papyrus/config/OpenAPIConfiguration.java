package org.purpurmc.papyrus.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.servers.Server;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
public class OpenAPIConfiguration {

    @Bean
    public OpenAPI getOpenAPI(AppConfiguration configuration) {
        return new OpenAPI()
                .info(new Info().title(configuration.getApiTitle()))
                .servers(List.of(new Server().url(configuration.getApiUrl())));
    }
}
