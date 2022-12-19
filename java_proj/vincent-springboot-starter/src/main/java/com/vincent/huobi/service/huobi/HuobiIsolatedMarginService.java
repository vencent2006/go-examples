package com.vincent.huobi.service.huobi;

import java.util.List;

import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;

import com.vincent.huobi.client.IsolatedMarginClient;
import com.vincent.huobi.client.req.margin.IsolatedMarginAccountRequest;
import com.vincent.huobi.client.req.margin.IsolatedMarginApplyLoanRequest;
import com.vincent.huobi.client.req.margin.IsolatedMarginLoanInfoRequest;
import com.vincent.huobi.client.req.margin.IsolatedMarginLoanOrdersRequest;
import com.vincent.huobi.client.req.margin.IsolatedMarginRepayLoanRequest;
import com.vincent.huobi.client.req.margin.IsolatedMarginTransferRequest;
import com.vincent.huobi.constant.Constants;
import com.vincent.huobi.constant.HuobiOptions;
import com.vincent.huobi.constant.Options;
import com.vincent.huobi.constant.enums.MarginTransferDirectionEnum;
import com.vincent.huobi.model.isolatedmargin.IsolatedMarginAccount;
import com.vincent.huobi.model.isolatedmargin.IsolatedMarginLoadOrder;
import com.vincent.huobi.model.isolatedmargin.IsolatedMarginSymbolInfo;
import com.vincent.huobi.service.huobi.connection.HuobiRestConnection;
import com.vincent.huobi.service.huobi.parser.isolatedmargin.IsolatedMarginAccountParser;
import com.vincent.huobi.service.huobi.parser.isolatedmargin.IsolatedMarginLoadOrderParser;
import com.vincent.huobi.service.huobi.parser.isolatedmargin.IsolatedMarginSymbolInfoParser;
import com.vincent.huobi.service.huobi.signature.UrlParamsBuilder;
import com.vincent.huobi.utils.InputChecker;

public class HuobiIsolatedMarginService implements IsolatedMarginClient {

  public static final String TRANSFER_TO_MARGIN_PATH = "/v1/dw/transfer-in/margin";
  public static final String TRANSFER_TO_SPOT_PATH = "/v1/dw/transfer-out/margin";
  public static final String GET_BALANCE_PATH = "/v1/margin/accounts/balance";
  public static final String GET_LOAN_ORDER_PATH = "/v1/margin/loan-orders";
  public static final String GET_LOAN_INFO_PATH = "/v1/margin/loan-info";

  public static final String APPLY_LOAN_PATH = "/v1/margin/orders";
  public static final String REPAY_LOAN_PATH = "/v1/margin/orders/{order-id}/repay";


  private Options options;

  private HuobiRestConnection restConnection;

  public HuobiIsolatedMarginService(Options options) {
    this.options = options;
    this.restConnection = new HuobiRestConnection(options);
  }


  @Override
  public Long transfer(IsolatedMarginTransferRequest request) {

    InputChecker.checker()
        .shouldNotNull(request.getDirection(), "direction")
        .checkSymbol(request.getSymbol())
        .checkCurrency(request.getCurrency())
        .shouldNotNull(request.getAmount(), "amount");

    String path = null;
    if (request.getDirection() == MarginTransferDirectionEnum.SPOT_TO_MARGIN) {
      path = TRANSFER_TO_MARGIN_PATH;
    } else {
      path = TRANSFER_TO_SPOT_PATH;
    }

    UrlParamsBuilder builder = UrlParamsBuilder.build()
        .putToPost("currency", request.getCurrency())
        .putToPost("symbol", request.getSymbol())
        .putToPost("amount", request.getAmount());

    JSONObject jsonObject = restConnection.executePostWithSignature(path, builder);
    return jsonObject.getLong("data");
  }

  @Override
  public Long applyLoan(IsolatedMarginApplyLoanRequest request) {

    InputChecker.checker()
        .checkSymbol(request.getSymbol())
        .checkCurrency(request.getCurrency())
        .shouldNotNull(request.getAmount(), "amount");

    UrlParamsBuilder builder = UrlParamsBuilder.build()
        .putToPost("currency", request.getCurrency())
        .putToPost("symbol", request.getSymbol())
        .putToPost("amount", request.getAmount());

    JSONObject jsonObject = restConnection.executePostWithSignature(APPLY_LOAN_PATH, builder);
    return jsonObject.getLong("data");
  }

  @Override
  public Long repayLoan(IsolatedMarginRepayLoanRequest request) {

    InputChecker.checker()
        .shouldNotNull(request.getOrderId(), "order-id")
        .shouldNotNull(request.getAmount(), "amount");

    UrlParamsBuilder builder = UrlParamsBuilder.build()
        .putToPost("amount", request.getAmount());

    String path = REPAY_LOAN_PATH.replace("{order-id}", request.getOrderId().toString());
    JSONObject jsonObject = restConnection.executePostWithSignature(path, builder);
    return jsonObject.getLong("data");
  }

  @Override
  public List<IsolatedMarginLoadOrder> getLoanOrders(IsolatedMarginLoanOrdersRequest request) {

    InputChecker.checker()
        .checkSymbol(request.getSymbol());

    UrlParamsBuilder builder = UrlParamsBuilder.build()
        .putToUrl("symbol", request.getSymbol())
        .putToUrl("start-date", request.getStartDate(), "yyyy-MM-dd")
        .putToUrl("end-date", request.getEndDate(), "yyyy-MM-dd")
        .putToUrl("states", request.getStatesString())
        .putToUrl("from", request.getFrom())
        .putToUrl("size", request.getSize())
        .putToUrl("direct", request.getDirection() == null ? null : request.getDirection().getCode());

    JSONObject jsonObject = restConnection.executeGetWithSignature(GET_LOAN_ORDER_PATH, builder);
    JSONArray data = jsonObject.getJSONArray("data");
    return new IsolatedMarginLoadOrderParser().parseArray(data);
  }

  @Override
  public List<IsolatedMarginAccount> getLoanBalance(IsolatedMarginAccountRequest request) {

    UrlParamsBuilder builder = UrlParamsBuilder.build()
        .putToUrl("symbol", request.getSymbol())
        .putToUrl("sub-uid", request.getSubUid());

    JSONObject jsonObject = restConnection.executeGetWithSignature(GET_BALANCE_PATH, builder);
    JSONArray data = jsonObject.getJSONArray("data");
    return new IsolatedMarginAccountParser().parseArray(data);
  }

  @Override
  public List<IsolatedMarginSymbolInfo> getLoanInfo(IsolatedMarginLoanInfoRequest request) {
    UrlParamsBuilder builder = UrlParamsBuilder.build()
        .putToUrl("symbols", request.getSymbols());

    JSONObject jsonObject = restConnection.executeGetWithSignature(GET_LOAN_INFO_PATH, builder);
    JSONArray data = jsonObject.getJSONArray("data");
    return new IsolatedMarginSymbolInfoParser().parseArray(data);
  }




}
