import Vue from 'vue'

// tslint:disable-next-line:match-default-export-name
import Vuex, { Store } from 'vuex'

import { Proj } from './modules/ProjStore'
import { IProjStoreState } from './modules/ProjStore/State'

Vue.use(Vuex)

const store: Store<IProjStoreState> = new Vuex.Store({
  modules: {
    Proj
  }
})

export default store
