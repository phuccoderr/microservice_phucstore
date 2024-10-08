package com.phuc.categoryservice.util;

import com.phuc.categoryservice.dtos.CategoryDto;
import com.phuc.categoryservice.exceptions.ParamValidateException;
import com.phuc.categoryservice.models.BaseEntity;
import com.phuc.categoryservice.models.Category;
import lombok.experimental.UtilityClass;
import org.modelmapper.ModelMapper;

import java.text.Normalizer;
import java.util.List;
import java.util.regex.Pattern;


@UtilityClass
public class Utility {

    static ModelMapper modelMapper = new ModelMapper();

    public String unAccent(String s) {
        String normalizer = Normalizer.normalize(s,Normalizer.Form.NFD);
        Pattern pattern = Pattern.compile("\\p{InCombiningDiacriticalMarks}+");
        String noAccents = pattern.matcher(normalizer).replaceAll("");
        return noAccents.toLowerCase().replace(" ", "-");
    }

    public List<CategoryDto> toListDtos(List<Category> categories) {
        return categories.stream().map(Utility::toDto).toList();
    }

    public CategoryDto toDto(Category category) {
        CategoryDto dto = modelMapper.map(category, CategoryDto.class);
        hierachical(dto);
        return dto;
    }

    public void hierachical(CategoryDto dto) {
        if (!dto.getChildren().isEmpty()) {
            dto.setHasChildren(true);
            dto.getChildren().forEach(item -> {
                hierachical(item);
            });
        } else {
            dto.setHasChildren(false);
        }
    }

    public void checkSortIsAscOrDesc(String sort) throws ParamValidateException {
        if (!sort.equals("asc") && !sort.equals("desc")) {
            throw new ParamValidateException("Param error sort: must be asc or desc");
        }
    }

}
