/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createTable('users', {
        id: 'id',
        name: { type: 'varchar(255)', notNull: true, unique: true },
        email: { type: 'varchar(255)', notNull: true, unique: true },
        password: { type: 'varchar(60)', notNull: true },
        createdAt: {
          type: 'timestamp',
          notNull: true,
          default: pgm.func('current_timestamp'),
        },
      })
    pgm.createTable('messages', {
        id: 'id',
        userId: {
          type: 'integer',
          notNull: true,
          references: '"users"',
          onDelete: 'cascade',
        },
        body: { type: 'text', notNull: true },
        createdAt: {
          type: 'timestamp',
          notNull: true,
          default: pgm.func('current_timestamp'),
        },
      })
      pgm.createIndex('messages', 'userId')
      pgm.createIndex('messages', 'createdAt')
};

exports.down = (pgm) => {
    pgm.dropTable('messages')
    pgm.dropTable('users')
};
