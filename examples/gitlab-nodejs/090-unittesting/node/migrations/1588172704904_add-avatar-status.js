/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
    pgm.createType('status_enum', ["Home", "Work", "Sleep"])
    pgm.addColumns('users', {
        avatarUrl: { type: 'varchar(60)'},
        status: { type: 'status_enum', default: 'Home' }
      })
};

exports.down = (pgm) => {
    pgm.dropColumns('users', ['avatarUrl', 'status'])
    pgm.dropType('status_enum')
};
