"""empty message

Revision ID: 1c09fc24ad9b
Revises: b34add8d70b0
Create Date: 2023-07-29 15:11:43.821324

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '1c09fc24ad9b'
down_revision = 'b34add8d70b0'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    with op.batch_alter_table('postcode', schema=None) as batch_op:
        batch_op.create_unique_constraint(None, ['postcode', 'locality'])

    with op.batch_alter_table('sa3', schema=None) as batch_op:
        batch_op.create_unique_constraint(None, ['code', 'name'])

    with op.batch_alter_table('sa4', schema=None) as batch_op:
        batch_op.create_unique_constraint(None, ['code', 'name'])

    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    with op.batch_alter_table('sa4', schema=None) as batch_op:
        batch_op.drop_constraint(None, type_='unique')

    with op.batch_alter_table('sa3', schema=None) as batch_op:
        batch_op.drop_constraint(None, type_='unique')

    with op.batch_alter_table('postcode', schema=None) as batch_op:
        batch_op.drop_constraint(None, type_='unique')

    # ### end Alembic commands ###
