import {
  Body,
  Controller,
  Get, HttpCode,
  HttpStatus,
  Param,
  Patch,
  Post,
  Query,
  Request,
  UploadedFile,
  UseGuards, UseInterceptors,
} from '@nestjs/common';
import { CustomersService } from './customers.service';
import { JwtAuthGuard } from '../auth/guards/jwt-auth.guard';
import { ResponseObject } from '../response/response-object.dto';
import { RedisCacheService } from '../redis/redis.service';
import { allCustomerKey, customerKey } from '../redis/key';
import { UpdateAccountDto } from './dto/update-account.dto';
import { ResponsePaginationDTO } from './dto/response-pagination.dto';
import { Customer } from './models/customer.schema';
import { RequestPaginationDto } from './dto/request-pagination.dto';
import { FileInterceptor } from '@nestjs/platform-express';
import { CloudinaryService } from '../cloudinary/cloudinary.service';
import { CUSTOMER_CONSTANTS } from '../constants/customer-constants';

@Controller('/api/v1/customers')
export class CustomersController {
  constructor(private readonly customersService: CustomersService,
              private readonly redisCacheService: RedisCacheService,
              private readonly cloudinaryService: CloudinaryService,) {}

  @UseGuards(JwtAuthGuard)
  @Get('account')
  async getAccount(@Request() req): Promise<ResponseObject> {
    const { _id } = req.user;

    const cacheCustomer = await this.redisCacheService.get(customerKey(_id));
    if (cacheCustomer) {
      return {
        data: cacheCustomer,
        status: HttpStatus.OK,
        message: CUSTOMER_CONSTANTS.GET_ACCOUNT,
      };
    }
    const customer = await this.customersService.findByCustomerId(_id);
    this.redisCacheService.set(customerKey(_id), customer);

    return {
      data: customer,
      status: HttpStatus.OK,
      message: CUSTOMER_CONSTANTS.GET_ACCOUNT,
    };
  }

  @UseGuards(JwtAuthGuard)
  @Patch('account')
  @HttpCode(HttpStatus.OK)
  async updateAccount(@Request() req,@Body() updateAccountDto: UpdateAccountDto): Promise<ResponseObject> {
    const { _id } = req.user;
    const customer = await this.customersService.updateCustomer(_id,updateAccountDto);
    this.redisCacheService.del(customerKey(_id));

    this.redisCacheService.set(_id, customer);

    return {
      data: customer,
      status: HttpStatus.OK,
      message: CUSTOMER_CONSTANTS.UPDATE_ACCOUNT,
    };
  }

  @UseGuards(JwtAuthGuard)
  @Get()
  async getAllCustomers(@Query() pagination: RequestPaginationDto): Promise<ResponseObject> {
    const { keyword, page, limit, sort } = pagination;

    if (!keyword) {
      const cachedAllUsers = await this.redisCacheService.get(
        allCustomerKey(page, limit, sort),
      );
      if (cachedAllUsers) {
        return {
          data: cachedAllUsers,
          status: HttpStatus.OK,
          message: CUSTOMER_CONSTANTS.GET_ALL,
        };
      }
    }

    const customers = await this.customersService.getAllCustomers(pagination)
    const paginationDto = this.buildPaginationDto(pagination, customers);

    customers.length !== 0 && this.redisCacheService.set(allCustomerKey(page, limit, sort), paginationDto);

    return {
      data: paginationDto,
      status: HttpStatus.OK,
      message: CUSTOMER_CONSTANTS.GET_ALL,
    };
  }

  @UseGuards(JwtAuthGuard)
  @Get(':id')
  async getUser(@Param('id') _id: string): Promise<ResponseObject> {
    const customer = await this.customersService.findByCustomerId(_id);

    return {
      data: customer,
      status: HttpStatus.OK,
      message: CUSTOMER_CONSTANTS.GET_ACCOUNT,
    };
  }

  @UseGuards(JwtAuthGuard)
  @UseInterceptors(FileInterceptor('file'))
  @Post('upload')
  @HttpCode(HttpStatus.OK)
  async uploadImage(@UploadedFile() file: Express.Multer.File,@Request() req): Promise<ResponseObject> {
    const { _id } = req.user;
    const customer = await this.customersService.findByCustomerId(_id);

    const result = await this.cloudinaryService.uploadFile(file);
    this.cloudinaryService.deleteImage(customer.image_id)

    await this.customersService.uploadAvatar(_id, result.url, result.public_id);

    this.redisCacheService.del(customerKey(_id));
    this.redisCacheService.set(_id, customer);

    return {
      data: {},
      status: HttpStatus.OK,
      message: CUSTOMER_CONSTANTS.UPLOAD_AVATAR,
    };
  }

  private buildPaginationDto(
    pagination: RequestPaginationDto,
    customers: Customer[],
  ): ResponsePaginationDTO {
    const { page, limit } = pagination;

    return {
      total_items: customers.length,
      total_pages: Math.ceil(customers.length / limit),
      current_page: parseInt(String(page)),
      start_count: (page - 1) * limit + 1,
      end_count: page * limit > customers.length ? customers.length : page * limit,
      entities: customers,
    };
  }

}
