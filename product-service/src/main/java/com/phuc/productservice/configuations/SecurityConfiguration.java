package com.phuc.productservice.configuations;

import com.phuc.productservice.constants.Constants;
import com.phuc.productservice.filter.JwtTokenFilter;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;
import org.springframework.web.cors.CorsConfiguration;
import org.springframework.web.cors.CorsConfigurationSource;
import org.springframework.web.cors.UrlBasedCorsConfigurationSource;

import java.util.List;

@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class SecurityConfiguration {

    private final JwtTokenFilter jwtTokenFilter;

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http.csrf(AbstractHttpConfigurer::disable)
                .authorizeHttpRequests( authz ->
                        authz
                                .requestMatchers(HttpMethod.GET, Constants.API_PRODUCTS).authenticated()
                                .requestMatchers(HttpMethod.GET, Constants.API_PRODUCTS + "/*").authenticated()
                                .requestMatchers(HttpMethod.POST, Constants.API_PRODUCTS).hasAnyAuthority("ADMIN")
                                .requestMatchers(HttpMethod.PATCH, Constants.API_PRODUCTS + "/*").hasAnyAuthority("ADMIN")
                                .requestMatchers(HttpMethod.DELETE, Constants.API_PRODUCTS + "/*").hasAnyAuthority("ADMIN")
                                .requestMatchers(HttpMethod.PATCH, Constants.API_PRODUCTS + "/add_files/*").hasAnyAuthority("ADMIN")
                                .requestMatchers(HttpMethod.DELETE, Constants.API_PRODUCTS + "/delete_files/*").hasAnyAuthority("ADMIN")
                                .anyRequest().permitAll())
                .sessionManagement(
                        session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS)
                )
                .addFilterBefore(jwtTokenFilter, UsernamePasswordAuthenticationFilter.class);
        return http.build();
    }

    @Bean
    CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration cc = new CorsConfiguration();
        cc.setAllowCredentials(true);
        cc.setAllowedOrigins(List.of("*"));
        cc.setAllowedHeaders(List.of("*"));
        cc.setAllowedMethods(List.of("*"));
        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/**", cc);
        return source;
    }
}
