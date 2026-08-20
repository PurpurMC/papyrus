package org.purpurmc.papyrus.config;

import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.filter.UrlHandlerFilter;

@Configuration
public class WebConfiguration {

    @Bean
    public FilterRegistrationBean<UrlHandlerFilter> trailingSlashFilter() {
        UrlHandlerFilter filter = UrlHandlerFilter
                .trailingSlashHandler("/**")
                .wrapRequest()
                .build();

        FilterRegistrationBean<UrlHandlerFilter> registrationBean = new FilterRegistrationBean<>(filter);
        registrationBean.setOrder(Integer.MIN_VALUE);
        return registrationBean;
    }
}
