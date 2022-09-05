// exmaple usage:
//  new Filters({
//    el     : '.metabar-filters',
//    query  : 'grizzly bears', /*current_user_searched_query*/
//    url : '/search/?q=$1' /*only supports '$1' as a match*/
//  });
//
// filters assumed in the form of <input type="checkbox" name="ABC" value="123"> which ouput "ABC:123"

(function(){

  var Filters = (function(){

    var filters = function(options){

      options.query = decodeURIComponent(options.query)

      if(typeof options.el !== 'undefined'){
        this.$el = $(options.el);
      }

      if(typeof options.url !== 'undefined'){
        this.url = options.url;
      }

      this.concat_char = "+";
      if(typeof options.concat_char !== 'undefined'){
        this.concat_char = options.concat_char;
      }

      this.param_char = ":";
      if(typeof options.param_char !== 'undefined'){
        this.param_char = options.param_char;
      }

      if(typeof options.query !== 'undefined'){
        // replaces existing params
        this.query = options.query.replace(/\w+:\w+/g,'');
      }


      this.filters = this.$el.find('[name]');

      this.bind();
      return this;
    };

    filters.prototype.parse = function(){

      var params  = [],
          checked = this.filters.filter('select:not(:disabled), :checked');

      for(var i=0,j=checked.length;i<j;i++){
        if (checked[i].name == "sort") {
          if (checked[i].value != "") {
            params.push( [checked[i].name,checked[i].value].join(this.param_char));
          }
        } else {
          params.push( [checked[i].name,checked[i].value].join(this.param_char));
        }
      }

      var hidden = this.filters.filter(':hidden');
      for(var i=0,j=hidden.length;i<j;i++){
        params.push( [hidden[i].name,hidden[i].value].join(this.param_char));
      }

      return params.join(this.concat_char);
    };

    filters.prototype.formSubmit = function(event) {
      event.preventDefault();
      this.query = $('#navbar-search-input').val();
      this.changeHandler();
    }

    filters.prototype.changeHandler = function(event) {

      var params          = this.parse(),
          new_query       = (this.query.length > 0 ? encodeURIComponent($.trim(this.query)) + "+" : encodeURIComponent($.trim(this.query))) + (params.length>0 ? params : ''),
          formatted_query = this.url.replace(/\$1/g, new_query);

     window.location.href = formatted_query;
    };

    filters.prototype.bind = function(){
      var events = [
        ['change', '[name]', 'changeHandler'],
      ];
      for(var i in events){
        this.$el.on(events[i][0], events[i][1], $.proxy(this[events[i][2]], this) );
      }
    };

    return filters;

  })();
  window.Filters = Filters;
})();