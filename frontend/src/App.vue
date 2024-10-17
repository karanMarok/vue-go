<template>
  <div>
    <button @click="makeApiCall">Click!</button>
    <pre>{{ message.message }}</pre>
  </div>
</template>


<script>
export default {
  data() {
    return {
      message: {}
    }
  },
  methods: {
    async makeApiCall() {
      try {
        const res = await fetch("http://localhost:8000/user", {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        })
        const newRes = await res.json();
        this.message = newRes

      } catch (error) {
        console.log(error)
      }
    }
  },
}
</script>