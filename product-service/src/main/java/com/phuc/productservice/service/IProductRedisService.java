package com.phuc.productservice.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.phuc.productservice.dtos.PaginationDto;
import com.phuc.productservice.exceptions.ParamValidateException;

public interface IProductRedisService {

    void clear();
    PaginationDto getAllCategories(Integer page, Integer limit, String sort) throws ParamValidateException, JsonProcessingException;
    void saveAllCategories(PaginationDto paginationDto, Integer page, Integer limit, String sort) throws JsonProcessingException;
}
