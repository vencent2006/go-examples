import { Module } from "@nestjs/common";
import { AppController } from "./app.controller";
import { AppService } from "./app.service";
import { AccountModule } from "./account/account.module";
import { TypeOrmModule } from "@nestjs/typeorm";
import { AuthModule } from "./auth/auth.module";

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: "mysql",
      host: "localhost",
      port: 3306,
      username: "root",
      password: "qwer1234",
      database: "nextseason",
      autoLoadEntities: true // 这个一定要写
    }),
    AccountModule,
    AuthModule],
  controllers: [AppController],
  providers: [AppService]
})
export class AppModule {
}
