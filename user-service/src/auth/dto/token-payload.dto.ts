import { ROLE } from "@src/auth/decorators/role.enum";

export class TokenPayLoad {
  _id: string;
  email: string;
  roles: ROLE[];
}
