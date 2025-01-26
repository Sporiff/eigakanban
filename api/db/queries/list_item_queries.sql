-- name: AddItemToListAtPosition :one
-- Arguments: list_uuid, item_uuid, new_position
WITH new_item AS (
    INSERT INTO list_items (list_id, item_id, position, prev_item_id, next_item_id)
        VALUES (
                   (SELECT list_id FROM lists WHERE lists.uuid = @list_uuid),
                   (SELECT item_id FROM items WHERE items.uuid = @item_uuid),
                   @new_position,
                   (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE uuid = @list_uuid) AND list_items.position = @new_position - 1),
                   (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE uuid = @list_uuid) AND list_items.position = @new_position + 1)
               )
        RETURNING list_item_id, list_id, position, prev_item_id, next_item_id
),
     update_prev_item AS (
         UPDATE list_items
             SET next_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.prev_item_id
                 AND new_item.prev_item_id IS NOT NULL
     ),
     update_next_item AS (
         UPDATE list_items
             SET prev_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.next_item_id
                 AND new_item.next_item_id IS NOT NULL
     )
SELECT list_items.uuid, list_items.position
FROM list_items
         JOIN items ON list_items.item_id = items.item_id
WHERE list_items.list_item_id = (SELECT list_item_id FROM new_item);

-- name: MoveItemInList :one
-- Arguments: list_item_uuid, new_position
WITH current_item AS (
    SELECT list_item_id, prev_item_id, next_item_id, position, list_id
    FROM list_items
    WHERE list_items.uuid = @list_item_uuid
),
     updated_positions AS (
         UPDATE list_items
             SET position = position + CASE
                                           WHEN list_items.position >= @new_position THEN 1
                                           ELSE -1
                 END
             WHERE list_id = (SELECT list_id FROM current_item)
                 AND position BETWEEN LEAST(@new_position, (SELECT position FROM current_item))
                       AND GREATEST(@new_position, (SELECT position FROM current_item))
             RETURNING list_item_id, position
     ),
     update_current_item AS (
         UPDATE list_items
             SET
                 position = @new_position,
                 prev_item_id = (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM current_item) AND list_items.position = @new_position - 1),
                 next_item_id = (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM current_item) AND list_items.position = @new_position + 1)
             WHERE list_item_id = (SELECT list_item_id FROM current_item)
             RETURNING uuid, position, prev_item_id, next_item_id
     )
SELECT uuid, position, prev_item_id, next_item_id
FROM update_current_item;

-- name: DeleteItemFromList :exec
-- Arguments: list_item_uuid
WITH deleted_item AS (
    DELETE FROM list_items
        WHERE list_items.uuid = @list_item_uuid
        RETURNING list_item_id, prev_item_id, next_item_id
),
     update_prev_item AS (
         UPDATE list_items
             SET next_item_id = deleted_item.next_item_id
             FROM deleted_item
             WHERE list_items.list_item_id = deleted_item.prev_item_id
                 AND deleted_item.prev_item_id IS NOT NULL
     ),
     update_next_item AS (
         UPDATE list_items
             SET prev_item_id = deleted_item.prev_item_id
             FROM deleted_item
             WHERE list_items.list_item_id = deleted_item.next_item_id
                 AND deleted_item.next_item_id IS NOT NULL
     )
SELECT 1;

-- name: MoveItemToAnotherList :one
-- Arguments: list_item_uuid, target_list_uuid, new_position
WITH current_item AS (
    SELECT list_item_id, item_id, prev_item_id, next_item_id, position, list_id
    FROM list_items
    WHERE list_items.uuid = @list_item_uuid
),
     update_source_list AS (
         UPDATE list_items
             SET
                 next_item_id = current_item.next_item_id
             FROM current_item
             WHERE list_items.list_item_id = current_item.prev_item_id
                 AND current_item.prev_item_id IS NOT NULL
     ),
     update_source_list_next AS (
         UPDATE list_items
             SET
                 prev_item_id = current_item.prev_item_id
             FROM current_item
             WHERE list_items.list_item_id = current_item.next_item_id
                 AND current_item.next_item_id IS NOT NULL
     ),
     new_item AS (
         INSERT INTO list_items (list_id, item_id, position, prev_item_id, next_item_id)
             VALUES (
                        (SELECT list_id FROM lists WHERE lists.uuid = @target_list_uuid),
                        (SELECT item_id FROM current_item),
                        @new_position,
                        (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE lists.uuid = @target_list_uuid) AND list_items.position = @new_position - 1),
                        (SELECT list_item_id FROM list_items WHERE list_id = (SELECT list_id FROM lists WHERE lists.uuid = @target_list_uuid) AND list_items.position = @new_position + 1)
                    )
             RETURNING list_item_id, list_id, position, prev_item_id, next_item_id
     ),
     update_destination_list AS (
         UPDATE list_items
             SET
                 next_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.prev_item_id
                 AND new_item.prev_item_id IS NOT NULL
     ),
     update_destination_list_next AS (
         UPDATE list_items
             SET
                 prev_item_id = new_item.list_item_id
             FROM new_item
             WHERE list_items.list_item_id = new_item.next_item_id
                 AND new_item.next_item_id IS NOT NULL
     )
SELECT uuid, position, prev_item_id, next_item_id
FROM list_items
WHERE list_items.list_item_id = (SELECT list_item_id FROM new_item);
