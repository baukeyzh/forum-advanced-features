$("#postTitle").on("input", function() {
  var input = document.getElementById('postTitle');
  //checks for deletion of text
  if (input.value.length > 30) {
    alert("input must have only 30 characters")
  // If it is, remove the last character from the input value
  input.value = input.value.slice(0, 30);
}
});
$("#postContent").on("input", function() {
  var input = document.getElementById('postContent');
  //checks for deletion of text
  if (input.value.length > 300) {
    alert("textarea must have only 300 characters")
  // If it is, remove the last character from the input value
  input.value = input.value.slice(0, 300);
}
});
$("input[type='file']").on("change", function () {
   if(this.files[0].size > 20000000) {
     alert("Please upload file less than 20MB. Thanks!!");
     $(this).val('');
   }
  });
  var isAuthorized = responseBody.autorized
  if (isAuthorized) {
    $("#authorized").addClass("hidden");
  } else {
    $(".logout").addClass("hidden");
  }
  posts = responseBody.posts
  $(".close").click(function() {
    $("#createPostModal").removeClass("hide");
});
$(document).click(function (e) {
  if ($(e.target).is('#createPostModal')) {
    $("#createPostModal").removeClass("hide");
  }

});
$(".post-create-btn").click(function() {
  if(!$("#createPostModal").hasClass("hide")) {
    $("#createPostModal").addClass("hide");
  } else {
    $("#createPostModal").removeClass("hide");
  }
});
$(".post-edit-btn").click(function() {
  alert($(this).data-id)
  if(!$("#createPostModal").hasClass("hide")) {
    $("#createPostModal").addClass("hide");
  } else {
    $("#createPostModal").removeClass("hide");
  }
});   
  
  
  function isSame(_arr1, _arr2) {
    if (
      !Array.isArray(_arr1)
      || !Array.isArray(_arr2)
      ) {
        return false;
      }
    
    // .concat() to not mutate arguments
    const arr1 = _arr1.concat().sort();
    const arr2 = _arr2.concat().sort();
    
    for (let i = 0; i < arr1.length; i++) {
      for (let j = 0; j < arr2.length; j++) {
        if (arr1[i] == arr2[j]) {
            return true;
         }
      }
    }
    
    return false;
}

  function getPosts() {
    var content = "<div>"
    if (posts) {
      console.log(posts)
      posts.forEach(function(post) {
        if (post.myLikeId != 0) {
          if (!post.categories.includes(", I like it")) {
            post.categories += ", I like it"
          }
        }
        if (post.authorId == responseBody.userId) {
          if (!post.categories.includes(", My post")) {
            post.categories += ", My post"
          }
        }
        var date = new Date(post.date)
        content += 
        `<div class="post">
          <div class="post-image">
            <img alt="" src="/imgs/` + post.imageName + `">
          </div>
          <div class="post-content">
            <div class="post-metadata">
              `+ getEditButton(post, responseBody.userId) +`
              <span class="post-date">` +date.toDateString() + `</span>
              <span class="post-category">` + post.categories + `</span>
              <span class="post-likes">
                <img style="width: 25px;" alt="" src="/imgs/thumbs-up-regular.svg">
                ` + post.likes + `
              </span>
              <span class="post-dislikes">
                <img style="width: 25px;" alt="" src="/imgs/thumbs-down-regular.svg">
                ` + post.dislikes +`
                </span>
            </div>
            <h2 class="post-title">
              <a href="/post-page?id=` + post.id +`">`+ post.title +`</a>
            </h2>
            <div class="post-text">`+ post.content +`</div>
            <a href="/post-page?id=` + post.id +`" class="read-more">Read More</a>
          </div>
        </div>
        `
      });
    }
    content  += '</div>';
    document.getElementById("allPosts").innerHTML = content;
  }
  

  $(document).ready(function() {
  getPosts()
  var categories = []
  $(".category-btn").click(function() {
    //$(".category-btn").removeClass("active");
   Array.prototype.remove = function() {
    var what, a = arguments, L = a.length, ax;
    while (L && this.length) {
        what = a[--L];
        while ((ax = this.indexOf(what)) !== -1) {
            this.splice(ax, 1);
        }
    }
    return this;
    };
    $(this).addClass("active");
    // Perform filtering based on selected category
    var selectedCategory = $(this).data("category");
    if (categories.indexOf(selectedCategory) == -1) {
      categories.push(selectedCategory);
    } else {
      $(this).removeClass("active")
      categories.remove(selectedCategory)
    }
    if (categories.length == 0 ) {
      getPosts()
    } else {
      newContent = "<div>"
      if (posts) {

        posts.forEach(function(post) { 
          postCategories = post.categories.split(', ')
          if (isSame(categories, postCategories)) {
            var date = new Date(post.date)
            newContent += 
            ` <div class="post">
                <div class="post-image">
                  <img alt="" src="/imgs/` + post.imageName + `">
                </div>
                <div class="post-content">
                  <div class="post-metadata">
                    `+ getEditButton(post, responseBody.userId) +`
                    <span class="post-date">` +date.toDateString() + `</span>
                    <span class="post-category">` + post.categories + `</span>
                    <span class="post-likes">
                      <img style="width: 25px;" alt="" src="/imgs/thumbs-up-regular.svg">
                      ` + post.likes + `
                    </span>
                    <span class="post-dislikes">
                      <img style="width: 25px;" alt="" src="/imgs/thumbs-down-regular.svg">
                      ` + post.dislikes +`
                    </span>
                  </div>
                  <h2 class="post-title">
                    <a href="/post-page?id=` + post.id +`">`+ post.title +`</a>
                  </h2>
                  <div class="post-text">`+ post.content +`</div>
                  <a href="/post-page?id=` + post.id +`" class="read-more">Read More</a>
                </div>
              </div>
            `
            }
          });
      }
      newContent += "</div>"
      document.getElementById("allPosts").innerHTML = newContent;
    }
  });
});
function getEditButton(post, userId) {
  console.log(post.authorId, userId, post.authorId == userId)
  if (post.authorId == userId) {
    return `
    <a href="/v1/post/edit?post-id=`+ post.id+`">
        <button class="post-edit" type="button">Edit Post</button>
    </a>
    <a href="/v1/post/delete?post-id=`+ post.id +`">
        <button class="post-edit" style="background-color: red" type="button"  onclick="confirmDelete(event)">Delete Post</button>
    </a>`
  }
  return ""
}
function confirmDelete(event) {
  if (confirm("Are you sure you want to delete this post?")) {
    // Proceed with the deletion
    window.location.href = event.target.parentNode.getAttribute("href");
  } else {
    // Cancel the deletion
    event.preventDefault();
  }
}