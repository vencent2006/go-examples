import { Injectable, UnauthorizedException } from "@nestjs/common";
import { JwtService } from "@nestjs/jwt";
import { AccountService } from "../account/account.service";
import { JWT_ACCESS_TOKEN_PERIOD, JWT_REFRESH_TOKEN_PERIOD } from "./auth.jwt.secret";
import { LoginAccountByMailDto, LoginAccountDto } from "../account/dto/login-account.dto";
import { ConfigService } from "@nestjs/config";

@Injectable()
export class AuthService {
  constructor(
    private accountService: AccountService,
    private jwtService: JwtService,
    private configService: ConfigService) {
  }

  async login(loginAccountDto: LoginAccountDto) {
    const account = await this.accountService.login(loginAccountDto);
    const access_token = this.getAccessToken(account.id, account.username);
    const refresh_token = this.getRefreshToken(account.id);
    return {
      access_token,
      refresh_token
    };
  }

  async loginByMail(loginAccountDto: LoginAccountByMailDto) {
    const account = await this.accountService.loginByMail(loginAccountDto);
    const access_token = this.getAccessToken(account.id, account.username);
    const refresh_token = this.getRefreshToken(account.id);
    return {
      access_token,
      refresh_token
    };
  }

  async refreshToken(oldToken: string) {
    // console.log(this.configService.get("email_user"));
    try {
      const data = this.jwtService.verify(oldToken);
      const account = await this.accountService.findOne(data.userid);
      const access_token = this.getAccessToken(account.id, account.username);
      const refresh_token = this.getRefreshToken(account.id);
      return {
        access_token,
        refresh_token
      };
    } catch (e) {
      throw new UnauthorizedException("token 已失效，请重新登录");
    }
  }

  getAccessToken(userId: number, username: string): string {
    return this.jwtService.sign({
      userId,
      username
    }, {
      expiresIn: JWT_ACCESS_TOKEN_PERIOD
    });
  }

  getRefreshToken(userId: number): string {
    return this.jwtService.sign({
      userId
    }, {
      expiresIn: JWT_REFRESH_TOKEN_PERIOD
    });
  }
}
