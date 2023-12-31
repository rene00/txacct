"""empty message

Revision ID: 6a44abef5afa
Revises: 77e93ace598d
Create Date: 2023-09-27 08:36:48.941540

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '6a44abef5afa'
down_revision = '77e93ace598d'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    with op.batch_alter_table('anzsic', schema=None) as batch_op:
        batch_op.drop_constraint('anzsic_description_key', type_='unique')

    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    with op.batch_alter_table('anzsic', schema=None) as batch_op:
        batch_op.create_unique_constraint('anzsic_description_key', ['description'])

    # ### end Alembic commands ###
