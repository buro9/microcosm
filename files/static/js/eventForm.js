//////////////
// location //
//////////////
(function(){


  /*
  *   location field toggle
  */
  var LOCATION_TYPE_AVAILABLE = 1,
      LOCATION_TYPE_TBA       = 0,
      location_type_options   = $('#location-options input[name=location-type]'),
      form_location           = $('.form-location');

  location_type = LOCATION_TYPE_AVAILABLE;

  function toggleLocation(){
    if (location_type === LOCATION_TYPE_AVAILABLE){
      form_location.show();
    }
    if (location_type === LOCATION_TYPE_TBA){
      form_location.hide();
      if (typeof formMap !== 'undefined'){
        formMap.reset();
        locationControl.find('textarea').val('');
      }
    }
  }

  location_type_options.on('change',function(e){
    var _this = $(e.currentTarget);
    location_type = parseInt(_this.val(),10);
    toggleLocation();
  });


  /*
  *   bind location textbox
  */
  var locationControl = $('#location-control');

  function geoQuery(locationQuery){
    var geoQueryURL = '/geocode/?q=';
    $.ajax({
      url   : geoQueryURL + locationQuery,
      type  : 'GET'
    }).success(geoQueryResult)
      .error(geoQueryError);
  }

  function geoQueryError(error){
    console.log(error);
  }

  function geoQueryResult(data, response, xhr){

    // And if it's not found, show an error
    if (!data.place_id) {
      return;
    }

    // Zoom to where it is
    // data.boundingbox = ["52.381053","52.6087058","-2.0336486","-1.7288578"]
    formMap.map.fitBounds(
      [
        [data.boundingbox[0], data.boundingbox[2]],
        [data.boundingbox[1], data.boundingbox[3]]
      ]
    );

    // And drop a pin
    var newLatLng = new L.LatLng(data.lat, data.lon);
    formMap
      .clearMarkers()
      .addMarker(newLatLng)
      .renderMarkers()
      .update(newLatLng)
      .save();
  }

  locationControl.on('click', '#locate', function(e){

    var parent        = $(e.delegateTarget),
        textarea      = parent.find('textarea'),
        locationQuery = textarea.val();

    if ($.trim(locationQuery)!==""){
      locationQuery = locationQuery.replace(/\n/g,' ');
      geoQuery(locationQuery);
    }

  });

})();


/////////////////////
//  people widget  //
////////////////////
(function(){

  'use strict'

  var subdomain = $('meta[name="subdomain"]').attr('content');

  var participating = new Participating({
    el         : '.list-participants',
    className  : 'list-people list-people-sm',
    static_url : subdomain
  });

  var peopleWidget = new PeopleWidget({
    el         : '#invite',
    is_textbox : true,
    static_url : subdomain,
    dataSource : subdomain + '/api/v1/profiles?disableBoiler&top=true&q=',
    invited    : typeof window.form_invitees !== 'undefined' &&
                        window.form_invitees !== "None" &&
                        window.form_invitees.length > 0 ? window.form_invitees : []
  });

  // update the hidden input box
  var invite_input_field = $('input[name="invite"]');
  var updateInvitedField = function(){
    invite_input_field.val(peopleWidget.invitedListToDelimitedString());
  };

  // update inviteObject field for form errors
  var inviteJSON_input_field = $('input[name="inviteObject"]');
  var updateInviteJSONField = function(obj){
    inviteJSON_input_field.val(JSON.stringify(obj));
  }


  // triggers when user clicks on a person in the autocomplete dropdown
  var render_particating = function(invited){
    if (invited.length > 0){
      participating.render(invited).show();
    }else{
      participating.hide();
    }
    peopleWidget.show();
    updateInvitedField();
    updateInviteJSONField(invited);
  }

  peopleWidget.onSelection(render_particating);

  // triggers when the user clicks on a person in the participants list
  participating.$el.on('click', '.remove', function(e){
    var id = e.currentTarget.rel,
        personToBeRemoved = {};

    for (var i=0,j=peopleWidget.people_invited.length;i<j;i++){
      if (peopleWidget.people_invited[i].id === id){
        personToBeRemoved = peopleWidget.people_invited[i];
        break;
      }
    }

    if (typeof personToBeRemoved.sticky === 'undefined'){
      peopleWidget.removePersonFromInvitedById(id).render();
      if (peopleWidget.people_invited.length>0){
        participating.render(peopleWidget.people_invited).show();
      }else{
        participating.hide();
      }
      updateInvitedField();
    }

  });

  // update on pageload
  if (peopleWidget.people_invited.length > 0){
    render_particating(peopleWidget.people_invited);
  }

})();

////////////////////
//  comment box  //
///////////////////
(function(){
  'use strict';
  if (typeof window.is_edit_form == 'undefined'){
    var replyBox = new simpleEditor({
      el : '.reply-box'
    });
  }
})();