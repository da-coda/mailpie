import config from '../config'
import SSE from "@/common/SSE";
const store = {
    name: 'mails',
    namespaced: true,
    state: {
        mails: []
    },
    getters: {

    },
    actions: {

    },
    mutations: {
        addMail(state, mail) {
            state.mails.push(mail)
        }
    }
}

//let sse = new SSE(config.sseUrl, store.commit("mails/addMail"))
let sse = new SSE(config.sseUrl, (data) => {
    console.log(data)
})
sse.run()
export default store