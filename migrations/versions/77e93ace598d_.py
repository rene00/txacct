"""empty message

Revision ID: 77e93ace598d
Revises: 6513b2f0b193
Create Date: 2023-09-26 20:29:26.416080

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '77e93ace598d'
down_revision = '6513b2f0b193'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table('organisation_source',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('name', sa.String(), nullable=False),
    sa.PrimaryKeyConstraint('id'),
    sa.UniqueConstraint('name')
    )
    with op.batch_alter_table('organisation', schema=None) as batch_op:
        batch_op.add_column(sa.Column('source_id', sa.Integer(), nullable=True))
        batch_op.add_column(sa.Column('organisation_source_id', sa.Integer(), nullable=False))
        batch_op.create_foreign_key(None, 'organisation_source', ['organisation_source_id'], ['id'])

    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    with op.batch_alter_table('organisation', schema=None) as batch_op:
        batch_op.drop_constraint(None, type_='foreignkey')
        batch_op.drop_column('organisation_source_id')
        batch_op.drop_column('source_id')

    op.drop_table('organisation_source')
    # ### end Alembic commands ###
