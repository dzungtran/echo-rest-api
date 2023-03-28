<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>DWorks - Simple Login</title>
    <script src="https://www.gstatic.com/firebasejs/8.6.8/firebase-app.js"></script>
    <script src="https://www.gstatic.com/firebasejs/8.6.8/firebase-analytics.js"></script>
    <!-- Firebase UI -->
    <script src="https://www.gstatic.com/firebasejs/ui/4.8.0/firebase-ui-auth.js"></script>
    <script src="https://www.gstatic.com/firebasejs/8.6.8/firebase-auth.js"></script>
    <link type="text/css" rel="stylesheet" href="https://www.gstatic.com/firebasejs/ui/4.8.0/firebase-ui-auth.css" />
    <script>
        // Your web app's Firebase configuration
        // For Firebase JS SDK v7.20.0 and later, measurementId is optional
        var firebaseConfig = JSON.parse({{ printf "%s" .config }});
        // Initialize Firebase
        firebase.initializeApp(firebaseConfig);
        // firebase.analytics();
        // firebase.auth().signInWithPopup(new firebase.auth.GoogleAuthProvider());
    </script>
    <script type="text/javascript">
      initApp = function() {
        firebase.auth().onAuthStateChanged(async function(user) {
          if (user) {
            console.log('id token...', await user.getIdToken())
            // User is signed in.
            var displayName = user.displayName;
            var email = user.email;
            var emailVerified = user.emailVerified;
            var photoURL = user.photoURL;
            var uid = user.uid;
            var phoneNumber = user.phoneNumber;
            var providerData = user.providerData;
            user.getIdToken().then(function(accessToken) {
              document.getElementById('sign-in-status').textContent = 'Signed in';
              document.getElementById('sign-in').textContent = accessToken;
              document.getElementById('account-details').textContent = JSON.stringify({
                displayName: displayName,
                email: email,
                emailVerified: emailVerified,
                phoneNumber: phoneNumber,
                photoURL: photoURL,
                uid: uid,
                accessToken: accessToken,
                providerData: providerData
              }, null, '  ');

              var xhttp = new XMLHttpRequest();
              xhttp.open("GET", "/me");
              xhttp.setRequestHeader("Authorization", "Bearer " + accessToken);
              xhttp.send();
              console.log(xhttp.responseText);
            });
          } else {
            // User is signed out.
            document.getElementById('sign-in-status').textContent = 'Signed out';
            document.getElementById('sign-in').textContent = 'Sign in';
            document.getElementById('account-details').textContent = 'null';
          }
        }, function(error) {
          console.log(error);
        });
      };

      window.addEventListener('load', function() {
        initApp()
      });
    </script>
  </head>
  <body>
    <h1>Welcome to My Awesome App</h1>
    <pre id="sign-in-status"></pre>
    <div id="sign-in"></div>
    <pre id="account-details"></pre>
  </body>
</html>