import { Inject, Injectable } from '@nestjs/common';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import { Cache } from 'cache-manager';

@Injectable()
export class RedisCacheService {
  constructor(@Inject(CACHE_MANAGER) private cacheManager: Cache) {}

  async get(key: string) {
    return await this.cacheManager.get(key)
  }

  async set(key: string, value: object) {
    await this.cacheManager.set(key, value);
  }

  async del(key: string) {
    await this.cacheManager.del(key);
  }

  async clearAllCustomerCache() {
    const keys = await this.cacheManager.store.keys('all_customers:*');

    if (keys.length > 0) {
      keys.forEach((key) => {
        this.cacheManager.store.del(key);
      });
    }
  }
}