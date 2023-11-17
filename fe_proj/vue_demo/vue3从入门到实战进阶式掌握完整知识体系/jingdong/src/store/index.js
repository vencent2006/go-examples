import { createStore } from 'vuex'

export default createStore({
  state: {
    // 购物车
    cartList: {
      // 第一层级：商铺id
      // 第二层级：商品id
    }
  },
  getters: {
  },
  mutations: {
    addItemToCart(state, payload) {
      const { shopId, productId, productInfo } = payload
      // console.log(shopId, productId, productInfo)
      // const cartList = state.cartList
      let shopInfo = state.cartList[shopId]
      if (!shopInfo) { shopInfo = {} } // shopInfo不存在就给一个初值
      let product = shopInfo[productId]
      if (!product) {
        console.log('shopInfo ', productId, ' not exist')
        product = productInfo
        product.count = 0
      }
      product.count += 1

      // 写回到state
      shopInfo[productId] = product
      state.cartList[shopId] = shopInfo

      console.log('shopInfo  ', shopInfo[productId])
    }
  },
  actions: {
  },
  modules: {
  }
})
