# Wiki Service

## Prefix

All routes to the API Layer begin with a version and the service.  
<br>
So, calls to the `wiki` service begin with: `/v1/wiki`  

---

## Routes

### HTTP `GET` Requests

| Type      | Route                                     | Arguments             | Description       |
| ---       | ---                                       | ---                   | ---               |
| `GET`     | `/pages{?ind=x&num=y}`                    | `ind`, `num`          | Returns a list of page info and content.  |
| `GET`     | `/pages/:id`                              | `:id`                 | Returns the info and content for the specified page. |
| `GET`     | `/pages/:id/revisions{?ind=x&num=y}`      | `:id`, `ind`, `num`   | Returns a list of the revisions for the specified page. |
| `GET`     | `/pages/:id/revisions/:rev`               | `:id`, `:rev`         | Returns the info and content for the specified revision of the specified page. |

#### Arguments
`ind`: the index to be the first item  
`num`: the number of entries to retrieve  
`:id`: the slug (or uuid) of the page  
`:rev`: the uuid of the page revision  
`{}`: content in curly braces is optional  
`x`, `y`: any integer  

---

### HTTP `POST` Requests

| Type      | Route                                     | Arguments             | Fields            |
| ---       | ---                                       | ---                   | ---               |
| `POST`    | `/pages/new`                              | N/A                   | `slug`, `name`, `author`, `archive_date`, `new_page`  |
| `POST`    | `/pages/:id/revisions`                    | `:id`                 | `page_id`, `author`, `new_page` |

#### `/pages/new`
**Description:** Creates a new page entry with the submitted info.  
This is implemented using a multipart form, with the fields being passed in as form data.  This is useful because it allows the `new_page` file to be passed in as a file, rather than just a string.
**Type:** `POST`
**Arguments:** None
**Fields:**
`slug`: the unique, human-readable identifier for the page  
    - all lowercase, kebab-case  
`name`: title of the page  
    - any string  
`author`: author identification  
    - not really implemented yet. using student id for now, but that definitely won't be the actual implementation.  
`archive_date`: date to mark page as archived/not relevant (if any)  
    - UTC+0 Time: `YYYY-MM-DD HH:MM:SS`  
    - optional; can be blank;  
`new_page`: the markdown file with the page content  

#### `/pages/:id/revisions`
**Description:** Creates a new revision of the specified page with the difference between the current content and that specified in the `new_page` file.  
This is implemented using a multipart form, with the fields being passed in as form data.  This is useful because it allows the `new_page` file to be passed in as a file, rather than just a string.  
**Type:** `POST`  
**Arguments:**  
`:id`: the slug (or uuid) of the page  
    - this isn't really used as of right now, as the id of the page is also passed in as a form field (`page_id`)  
**Fields**:  
`page_id`: the slug (or uuid) of the page being revised  
`author`: author identification  
    - not really implemented yet. using student id for now, but that definitely won't be the actual implementation.  
`new_page`: the markdown file with the new page content  

