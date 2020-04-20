declare global {
  namespace NodeJS {
    export interface ProcessEnv {
      NODE_ENV: 'development' | 'production' | 'test'
      REACT_APP_INFURA_KEY: string
      REACT_APP_GA_ID?: string
      REACT_APP_FEEDS_JSON?: string
      REACT_APP_NODES_JSON?: string
    }
  }
}
