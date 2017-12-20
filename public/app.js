new Vue({
  el: '#app',

  data: {
      ws: null, // websocket
      word:[],
      letter: null
  },

  methods: {
    send: function(){
      if( !this.letter.match(/[A-Za-z]{1}/) ){
        alert("INVALID input");
        this.letter = "";
      }

      this.ws.send(JSON.stringify({
        type: "turn",
        data: this.letter
      }));

      this.letter = "";
    }
  },

  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');

    this.ws.onmessage = function(event){
      self.word = JSON.parse(event.data).word;
    };

  }

});
