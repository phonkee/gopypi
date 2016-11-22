<template>
    <div class="container">
        <div class="row">
            <div class="col-md-4 col-md-offset-4">
                <div class="login-panel panel panel-default">
                    <div class="panel-heading">
                        <h3 class="panel-title">Please Sign Into Gopypi</h3>
                    </div>
                    <div class="panel-body">
                        <form role="form" v-on:submit.prevent="doLogin">
                            <fieldset>
                                <div class="form-group">
                                    <input class="form-control" placeholder="Username" autofocus v-model="credentials.username">
                                </div>
                                <div class="form-group">
                                    <input class="form-control" placeholder="Password" type="password" v-model="credentials.password">
                                </div>
                                <a class="btn btn-lg btn-success btn-block" @click="doLogin()">Login</a>
                            </fieldset>
                        </form>
                    </div>
                    <div class="panel-footer">
                        <messages></messages>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import auth from '../api/auth'
import Messages from './Messages.vue'

export default {
  components: {
    Messages
  },
  data () {
    return {
      credentials: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    doLogin () {
      auth.login(this.credentials.username, this.credentials.password).then((token) => {
        this.$router.push({name: 'admin.dashboard'})
      }, (response) => {
        if (response.status !== 500 && response.data.error) {
          this.$store.dispatch('messageError', response.data.error)
        }
      })
    }
  }
}
</script>
