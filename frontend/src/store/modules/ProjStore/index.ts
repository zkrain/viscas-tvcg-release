import { Module } from 'vuex'

import { actions } from './Actions'
import { getters } from './Getters'
import { mutations } from './Mutations'
import { IProjStoreState, state } from './State'

import { IRootState } from '../../RootState'

// tslint:disable:variable-name
export const Proj: Module<IProjStoreState, IRootState> = {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
